#!/bin/bash
set -e

# Create S3 bucket for image storage
awslocal s3 mb s3://fitgenie-images

# Set bucket CORS policy for image uploads
awslocal s3api put-bucket-cors --bucket fitgenie-images --cors-configuration '{
    "CORSRules": [
        {
            "AllowedHeaders": ["*"],
            "AllowedMethods": ["PUT", "POST", "GET"],
            "AllowedOrigins": ["*"],
            "MaxAgeSeconds": 3000
        }
    ]
}'

echo "S3 bucket 'fitgenie-images' created and configured"
