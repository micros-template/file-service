package constant

const (
	APP_BUCKET           = "dropboks-bucket"
	PROFILE_IMAGE_FOLDER = "profile"
	PUBLIC_PERMISSION    = `{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {"AWS": "*"},
        "Action": ["s3:GetObject"],
        "Resource": ["arn:aws:s3:::%s/*"]
      }
    ]
  }`
)

const MAX_IMAGE_SIZE_BYTES = 6 * 1024 * 1024
