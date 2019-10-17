# Image Usecase

The usecase for image is only 2(at least for now):

1. Upload Image
2. Download Image

## Upload Image

There are two `mode` for  uploading image:

1. Public
2. Private

We want to have different mode, because not all image is public. We might want to store a sensitive image related to customer.

### Private Image

When uploading the `private` image, `metadata` attribute is saved alongside the `image`. In this `metadata`, contains the access and priviledge that belong to spesific user. This to make sure that `private` image can only be accessed by the rightful users.

## Download Image

Download image usecase exists to serve `private` image. Image that belong to one user and is `private` should not visible to other users. So we need to check whether the downloader is the rightful user.

### Temporary Path

To be added

### Object Storage Signed URL

To be added

## TODO

- Image manipulation(resize)
- Image compression
