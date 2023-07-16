# Disk Cleaning Utility

This Disk Cleaning Utility is a powerful and flexible tool for helping users reclaim storage space on their machines. It can identify and remove unnecessary `node_modules` directories and duplicate files that may be taking up significant space.

## Features

- **Delete `node_modules` directories**: Many JavaScript projects use `node_modules` directories to store project dependencies. Over time, these can accumulate and take up substantial disk space. This feature can find and delete all `node_modules` directories in the specified path.
- **Delete duplicate files**: Duplicate files can accumulate over time and waste storage space. This feature can identify duplicates based on file hash comparison and prompt the user for deletion.
- **Interactive prompts**: The utility asks users for their preferences for deletion, providing a personalized cleaning experience.
- **Progress display**: The utility provides a real-time update on how much disk space has been freed and the time elapsed since the operation started, giving users insight into the cleaning process.
  
## Usage

To run the program, make sure you have Go installed and your environment set up properly. Then follow these steps:

1. Clone this repository to your local machine.
2. Navigate to the project's directory.
3. Run `go build` to compile the program.
4. Run the executable file produced by the build command.

Once you start the program, you'll be prompted to enter the path to be cleaned and whether you'd like to delete `node_modules` directories and duplicate files. You can also specify a minimum file size for duplicate deletion.

## Note

Please be careful when using this or any disk cleaning utility. Always make sure your important data is backed up, and double-check any files before you agree to delete them.

## License

This project is licensed under the terms of the MIT license. See the [LICENSE](LICENSE) file for details.

