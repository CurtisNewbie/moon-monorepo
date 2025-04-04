# Maintenance

mini-fstore automatically detects duplicate uploads by comparing size and sha1 checksum. If duplicate file is detected, these files are _symbolically_ linked to the same file previously uploaded. This can massively reduce file storage, but multiple file records (multiple file_ids) can all point to a single file.

Whenever a file is marked logically deleted, the file is not truely deleted. In order to cleanup the storage for the deleted files including those that are possibly symbolically linked, you have to use the following endpoint to trigger the maintenance process. During the maintenance, uploading files is rejected.

```sh
curl -X POST http://localhost:8084/maintenance/remove-deleted
```

mini-fstore also provides maintenance endpoint that sanitize storage directory. Sometimes files are uploaded to storage directory, but are somehow not saved in database. These <i>dangling</i> files are handled by this endpoint.

```sh
curl -X POST http://localhost:8084/maintenance/sanitize-storage
```

To compute sha1 for previously uploaded files, use the following maintenance endpoint to trigger a compensation.

```sh
curl -X POST 'http://localhost:8084/maintenance/compute-checksum'