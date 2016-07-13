# Kubernetes Secret Env

Take Kubernetes secrets provided via a mounted volume and execute a process that has environment variables populated from secrets.

`kubernetes-secret-env {{ your program }}`

## Releasing new versions

To make the compiled version of `kubernetes-secret-env` available for Dockerfiles download, we have to separately attach that compiled file to the release on GitHub. This file has to be compiled on the same system architecture that you want it to run on.

On that system, run:

```bash
# Install golang
apt-get update && apt-get install -y golang vim

# Copy the `kubernetes-secret-env.go` source code to the system
vim kubernetes-secret-env.go # copy paste

go build
```

You then need to download the compiled file back to your system to upload it to GitHub. This will depend on what your remote system is.

If it happens to be a Docker container running on Google Kubernetes Engine:
```bash
# Locally
gsutil signurl -p notasecret -c "application/octet-stream" -m PUT [PATH TO PRIVATE KEY] gs://[GCS BUCKET]/kubernetes-secret-env

# Remotely
curl -XPUT -H "Content-Type: application/octet-stream" --data-binary @kubernetes-secret-env "[URL FROM ABOVE]"

# Locally
gsutil cp gs://[GCS BUCKET]/kubernetes-secret-env .
```

And you've got your file!

## Changelog

* `0.0.2` - Fixes #1
* `0.0.1` - Initial release
