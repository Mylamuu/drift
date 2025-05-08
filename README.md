# Drift

Drift a simple temporary file uploading service for transferring files between machines on my local network. Each uploaded file is assigned a unique ID and is stored on disk until the keep time is reached (by default 24 hours.) At any point during this time, anyone can download the file using the unique ID.

## Building

```sh
git clone https://github.com/Mylamuu/drift.git
cd drift && sh ./build.sh
```

The compile binary can then be found in the `bin` folder.

## License

This project is licensed under the [AGPL-3.0 License](LICENSE).