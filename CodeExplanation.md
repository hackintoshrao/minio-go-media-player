## Code Explanation

List your objects.

[ListObjects](https://github.com/minio/minio-go/blob/master/examples/s3/listobjects.go) is called on the specified bucket when player is initialized. These objects will be rendered as playlist for media player as shown in the player image above.

```golang
for objectInfo := range api.storageClient.ListObjects(*bucketName, "", isRecursive, doneCh) {
  if objectInfo.Err != nil {
		http.Error(w, objectInfo.Err.Error(), http.StatusInternalServerError)
			return
		}
		objectName := objectInfo.Key // object name.
		playListEntry := mediaPlayList{
			Key: objectName,
		}
		playListEntries = append(playListEntries, playListEntry)
	}
	playListEntriesJSON, err := json.Marshal(playListEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Successfully wrote play list in json.
	w.Write(playListEntriesJSON)

```
"Secure URLs are generated on demand when requested to play, underlying mechanism is [PreSignedGetObject]((https://github.com/minio/minio-go/blob/master/examples/s3/presignedgetobject.go)). This secure URL is used by the player to stream and play the media from the bucket."
```golang
// GetPresignedURLHandler - generates presigned access URL for an object.
func (api mediaHandlers) GetPresignedURLHandler(w http.ResponseWriter, r *http.Request) {
	// The object for which the presigned URL has to be generated is sent as a query
	// parameter from the client.
	objectName := r.URL.Query().Get("objName")
	if objectName == "" {
		http.Error(w, "No object name set, invalid request.", http.StatusBadRequest)
		return
	}
	presignedURL, err := api.storageClient.PresignedGetObject(*bucketName, objectName, 1000*time.Second, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(presignedURL))
}
```
