package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	httprequest "bitbucket.org/ayopop/of-core/http/v2"
)

type Data struct {
	CorrelationId              string
	TransactionReferenceNumber string
	BeneficiaryId              string
	CustomerId                 string
}

func Generate32Characters() string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var length int = 32
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	filePath := "data.csv" // Replace with your CSV file path

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Create an array to store the data
	var dataArray []Data

	// Iterate over the records and populate the dataArray
	for _, record := range records {
		if len(record) >= 4 {
			data := Data{
				CorrelationId:              record[0],
				TransactionReferenceNumber: record[1],
				CustomerId:                 record[2],
				BeneficiaryId:              record[3],
			}
			dataArray = append(dataArray, data)
		}
	}

	// Print the data array

	fmt.Println(dataArray)
	for _, v := range dataArray {
		HitStatus(v)
	}
}
func HitStatus(data Data) {

	clientSetting := &http.Client{}

	httpClient := httprequest.NewClient(clientSetting)
	rc := httprequest.NewRestClient(httpClient)
	var resp interface{}
	q := url.Values{}

	q.Set("transactionId", Generate32Characters())
	// q.Set(constant.MerchantCodeString, req.MerchantCode)
	q.Set("beneficiaryId", data.BeneficiaryId)
	q.Set("customerId", data.CustomerId)
	q.Set("transactionReferenceNumber", data.TransactionReferenceNumber)
	uri := fmt.Sprintf("https://api.of.ayoconnect.id/api/v1/bank-disbursements/status/%v", data.CorrelationId)

	httpResponse := rc.NewRequest().
		AddHeaders("Content-Type", "application/json").
		AddHeaders("Accept", "application/json").
		AddHeaders("A-Latitude", "234.8503").
		AddHeaders("A-Correlation-ID", Generate32Characters()).
		AddHeaders("A-Longitude", "23.654").
		AddHeaders("A-Merchant-Code", "AYOPOP").
		AddHeaders("A-IP-ADDRESS", "192.168.1.16").
		AddHeaders("productCode", "BULK_DISBURSEMENT").
		AddHeaders("Authorization", "Bearer hgA93ueGgwsrAsoaFK1ggrQkEI44").
		WithQuery(&q).
		Get(uri)

	respClientErr := httpResponse.
		Json(&resp).
		Error()
	if respClientErr != nil {
		fmt.Println(respClientErr)
	}
	fmt.Println(resp)
}
