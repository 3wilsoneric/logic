package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const fastAPIURL = "http://your-fastapi-backend/custom_predict"  // Replace with your actual FastAPI endpoint

func sendFastAPIRequest(inputData string) (string, error) {
	data := map[string]string{"text": inputData}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fastAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func handlePredictRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		inputData := r.FormValue("input_data")
		result, err := sendFastAPIRequest(inputData)
		if err != nil {
			http.Error(w, "Error sending request to FastAPI", http.StatusInternalServerError)
			return
		}

		// Process the result as needed
		fmt.Fprintf(w, "Prediction Result: %s", result)
	} else {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/predict", handlePredictRequest)
	http.ListenAndServe(":8080", nil)
}
