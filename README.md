# waifu2x-web

A worker based upscaler utilizing waifu2x.

## Workers
You need to install Beanstalkd v1.12 as the message brokers for the workers, and set up shared folder for the worker to process the images.

I'm gonna use Arch Linux as the example here.
1. Clone this repo
```
git clone https://github.com/hiyorun/waifu2x-upscaler-web
```

2. Compile

worker-vulkan and waifu2x-ncnn-vulkan:
```
pacman -S waifu2x-ncnn-vulkan
cd waifu2x-upscaler-web
cd worker-vulkan
go build
```

or worker-cpp and waifu2x-converter-cpp:
```
yay -S waifu2x-converter-cpp
cd waifu2x-upscaler-web
cd worker-cpp
go build
```

Note: You can install both on the same machine. Deploy them to your instance(s). Use `--help` for connection details to beanstalkd and main backend

## Backend
```
cd waifu2x-upscaler-web
go build
```

Deploy it to the central instance that saves the processed images to an SQLite DB (See `--help`)

## Frontend
```
cd waifu2x-upscaler-web
yarn
yarn build
```

Currently the frontend will use current domain or localhost. No configuration is possible currently unless you edit the useAPI composables. Sorry!
