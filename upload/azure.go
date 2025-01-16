package upload

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"lab.sda1.net/nexryai/hoyofeed/logger"
	"os"
	"strings"
)

var (
	connectionString = os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	containerName    = os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
)

func PutFileToAzureBlob(filePath, blobName string) error {
	log := logger.GetLogger("UPLD")

	contentType := "application/octet-stream"
	if strings.HasSuffix(filePath, ".json") {
		contentType = "application/json"
	} else if strings.HasSuffix(filePath, ".xml") {
		// RSS feed
		contentType = "application/xml"
	}

	client, err := azblob.NewClientFromConnectionString(connectionString, nil)

	// Load the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	log.Info("Uploading file to Azure Blob Storage...")

	_, err = client.UploadFile(context.TODO(), containerName, blobName, file, &azblob.UploadFileOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr(contentType),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	} else {
		log.Info(fmt.Sprintf("File uploaded successfully to container %s with blob name %s", containerName, blobName))
	}

	return nil
}
