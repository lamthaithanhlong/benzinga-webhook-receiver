package config

import (
	"benzinga-webhook-receiver/pkg/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var ticker *time.Ticker

var (
	batch []models.LogEntry
	mutex sync.Mutex
)

type LogEntry struct {
	UserID    int                    `json:"user_id"`
	Total     float64                `json:"total"`
	Title     string                 `json:"title"`
	Meta      map[string]interface{} `json:"meta"`
	Completed bool                   `json:"completed"`
}

type Config struct {
	Port          string        `yaml:"port"`
	BatchSize     int           `yaml:"batch_size"`
	BatchInterval time.Duration `yaml:"batch_interval"`
	PostEndpoint  string        `yaml:"post_endpoint"`
}

var Cfg Config

func LoadConfig(path string) {
	logrus.Infof("Loaded configuration: %+v", Cfg)
	yamlFile, err := os.ReadFile(path) // Updated from ioutil.ReadFile to os.ReadFile
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
	}
	err = yaml.Unmarshal(yamlFile, &Cfg)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %s\n", err)
	}

	// Override YAML settings with environment variables if they exist
	if port := os.Getenv("PORT"); port != "" {
		Cfg.Port = port
	}
	if batchSize := os.Getenv("BATCH_SIZE"); batchSize != "" {
		Cfg.BatchSize, _ = strconv.Atoi(batchSize)
	}
	if batchInterval := os.Getenv("BATCH_INTERVAL"); batchInterval != "" {
		Cfg.BatchInterval, _ = time.ParseDuration(batchInterval)
	}
	if postEndpoint := os.Getenv("POST_ENDPOINT"); postEndpoint != "" {
		Cfg.PostEndpoint = postEndpoint
	}
}

func AddToBatch(entry models.LogEntry) {
	mutex.Lock()
	batch = append(batch, entry)
	mutex.Unlock()

	if len(batch) >= Cfg.BatchSize {
		sendBatch()
	}
}

func sendBatch() {
	logrus.Infof("Attempting to send batch to endpoint: %s", Cfg.PostEndpoint)
	mutex.Lock()
	if len(batch) == 0 {
		mutex.Unlock()
		return
	}

	// Prepare the payload
	payload, err := json.Marshal(batch)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal batch data")
		mutex.Unlock()
		return
	}

	attempts := 0
	success := false

	for attempts < 3 && !success {
		req, err := http.NewRequest("POST", Cfg.PostEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			logrus.WithError(err).Error("Failed to create request")
			time.Sleep(2 * time.Second)
			attempts++
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.WithError(err).Error("Failed to send batch")
			time.Sleep(2 * time.Second)
			attempts++
		} else {
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				logrus.Infof("Successfully sent batch of %d entries", len(batch))
				success = true
			} else {
				logrus.WithField("status_code", resp.StatusCode).Error("Failed to send batch")
				time.Sleep(2 * time.Second)
				attempts++
			}
		}
	}

	if !success {
		logrus.Fatal("Failed to send batch after 3 attempts, exiting application")
		os.Exit(1)
	}

	// Clear the batch if successfully sent
	if success {
		batch = []models.LogEntry{} // Clear the batch if successful
	}
	mutex.Unlock()
}

func StartBatchTimer() {
	if Cfg.BatchInterval <= 0 {
		logrus.Error("Invalid batch interval in configuration; setting to default 5m")
		Cfg.BatchInterval = 5 * time.Minute
	}
	ticker = time.NewTicker(Cfg.BatchInterval)
	go func() {
		for range ticker.C {
			sendBatch()
		}
	}()
}

func StopBatchTimer() {
	if ticker != nil {
		ticker.Stop()
		ticker = nil
	}
}
