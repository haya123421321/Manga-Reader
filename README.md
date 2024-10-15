# Manga Viewer Web Application

Welcome to the Manga Viewer web application! This project represents my first deep dive into web application development, and I have gained invaluable insights and skills throughout this journey.

## Overview

The Manga Viewer is a simple yet functional web application that allows users to browse and read manga offline. It features a user-friendly interface where you can navigate through various chapters and read your favorite manga. Additionally, it includes a bookmarking feature that remembers your last read chapter, so you never have to worry about losing your place.
![Home showcase](https://github.com/user-attachments/assets/d013086a-45cc-4704-8b90-9362c1c62619)
![Reading showcase](https://github.com/user-attachments/assets/e7ee0313-6bd1-4b81-999c-aa60b9431435)


### Features

- **Chapter Navigation**: Easily move between previous and next chapters.
- **Image Display**: View manga pages in a clean layout, optimized for all screen sizes.
- **Dynamic Content**: The application reads from a local directory structure to serve manga chapters and images.
- **Cross Platform**: The biggest reason why I went the web app route is so more people can enjoy this experience seamlessly, regardless of their device or operating system

## Technology Stack

This project is built using:

- **Go**: The backend is powered by Go, utilizing its `net/http` package to handle HTTP requests and serve content.
- **HTML/CSS**: The front-end interface is designed with HTML and CSS for structure and styling.

## Installation

### Windows
Visit the [Releases](https://github.com/haya123421321/Manga-Reader/releases) section of this repository to download the latest compiled binary for Windows.
Extract the Files

Once downloaded, extract the contents of the ZIP file to a folder of your choice.

### Linux
   ```bash
   git clone https://github.com/haya123421321/Manga-Reader.git
   cd Manga-Reader
   go build -o Manga-Reader main.go
   ```

## Post
Place your manga files in the `mangas` directory, like so:
```
├── mangas
│   └── Berserk
│       └── Berserk 001
│           ├── 0.jpg
│           ├── 1.jpg
│           ├── 2.jpg
│           ├── 3.jpg
│           ├── 4.jpg
│           ├── ...
└── README.md

```


Start the program and  navigate to `http://127.0.0.1:8080` in your browser to access the application.

## Future Improvements

As this is my first web application, I'm still learning and planning to implement new features, including:

- Improved error handling and user feedback.
- Enhanced styling for a better user experience.
- Add support for zip and cbz files

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Feel free to contribute to this project! Open an issue or submit a pull request if you have suggestions or improvements.

And lasty thank you to ChatGPT for writing this README, couldn't have done it without you.
