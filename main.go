package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type manga_page_data struct {
	Images []string
	Name string
	Next string
	Previous string
	Chapters_names []string
	Chapters_links []string
}

type manga_information struct {
	Title string `json:title`
	Data reading_data `json:data`
}

type reading_data struct {
	Last_read_chapter string `json:last_read_chapter`
}

var (
	Titles []string
	Icons []string
	handlers_made []string
) 


func main() {
	http.Handle("/mangas/", http.StripPrefix("/mangas/", http.FileServer(http.Dir("./mangas"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", home)

	directoryss,err := os.ReadDir("mangas")
	if err != nil {
		os.Mkdir("mangas", 0700)
	}

	for _,directory := range directoryss {
		if directory.IsDir() {
			Titles = append(Titles, directory.Name())
			url_encode_title := url.PathEscape(directory.Name())
			http.HandleFunc("/" + url_encode_title + "/", func(w http.ResponseWriter, r *http.Request) {
				manga_page(w, r, directory.Name())
			})

		Subdir,_ := os.ReadDir("mangas/" + directory.Name())
		first_chapter,err := os.ReadDir("mangas/" + directory.Name() + "/" +  Subdir[0].Name()) 

		if err != nil {
			fmt.Println("Something went wrong when trying to read file:")
			fmt.Println(Subdir[0].Name())
		}
		
		first_image := "/mangas/" + directory.Name() + "/" + Subdir[0].Name() + "/" + first_chapter[0].Name()
		Icons = append(Icons, first_image)
		}
	}
	
	var port string

	for i,arg := range os.Args {
		if arg == "-p" {
			port = os.Args[i + 1]
		}
	}

	if port != "" {
		fmt.Println("Listening on port " + port)
		fmt.Println("Go to http://127.0.0.1:" + port + " to access the viewer.")
		err = http.ListenAndServe("0.0.0.0:" + port, nil)
	} else {
		fmt.Println("Listening on port 8080")
		fmt.Println("Go to http://127.0.0.1:8080 to access the viewer.")
		err = http.ListenAndServe("0.0.0.0:8080", nil)
	}
}

type home_data struct {
	Titles []string
	Icons []string
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl,err := template.ParseFiles(filepath.Join("template", "home.html"))
	if err != nil {
		fmt.Println("Coudlnt parse home.html")
		fmt.Println(err)
		return
	}

	d := home_data {
		Titles: Titles,
		Icons: Icons,
	}

	tmpl.Execute(w, d)
	
}

func manga_page(w http.ResponseWriter, r *http.Request, title string) {
	file,err := os.ReadFile("data.json")
	if err != nil {
		os.Create("data.json")
	}
	
	var json_data []manga_information
	json.Unmarshal(file, &json_data)
	
	var title_in_data bool
	for _,j := range json_data {
		if j.Title == title {
			title_in_data = true
			break
		}
	}
	
	manga_path := "/mangas/" + title
	files,err := os.ReadDir("." + manga_path)
	if err != nil {
		fmt.Println("Couldnt read from manga directory " + manga_path)
	}


	var chapters []string
	for _,file := range files {
		if file.IsDir() {
			chapters = append(chapters, file.Name())
		}
	}
	
	var manga_index_in_json int
	var chapter_index = 0
	if !title_in_data {
		data := reading_data{Last_read_chapter: chapters[0]}
		json_data = append(json_data, manga_information{Title: title, Data: data})
		new_json_data,_ := json.MarshalIndent(json_data, "", "\t")
		os.WriteFile("data.json", new_json_data, 0600)
		
		chapter_index = 0
	} else {
		manga_index_in_json = findMangaIndexFromJson(json_data, title)
		for i,chapter := range chapters {
			if json_data[manga_index_in_json].Data.Last_read_chapter == chapter {
				chapter_index = i
				break
			}
		}
	}


	chapter := chapters[chapter_index]

	
	for _,c := range chapters {
		safeHandler("/" + url.PathEscape(title) + "/" + url.PathEscape(c), func(w http.ResponseWriter, r *http.Request) {
			load_images(w, title, chapters, manga_path, c)
		})
	}


	http.Redirect(w, r, "/" + title + "/" + chapter, http.StatusSeeOther)
}

func load_images(w http.ResponseWriter, title string, chapters []string, manga_path string, chapter string) {
	chapter_path := manga_path + "/" + chapter
	image_files,err := os.ReadDir("." + chapter_path)
	if err != nil {
		return
	}
	
	var files_string []string
	
	// Custom sorting through image files
	if len(image_files) > 9 {
		var Not_Numeric_Format_Checker int
		var Is_Numeric_Format_Checker int

		for _,file := range image_files {
			if Not_Numeric_Format_Checker == 3 {
				fmt.Println("Images not using numeric formatting, using default sort")
				files_string = []string{}
				for _,file := range image_files {
					files_string = append(files_string, chapter_path + "/" + file.Name())
				}
				break
			} else if Is_Numeric_Format_Checker == 3 {
				files_string = []string{}
				for _,i := range image_files {
					files_string = append(files_string, chapter_path + "/" + i.Name())
				}

				sort.Slice(files_string, func(i, j int) bool {
					split_i := strings.Split(files_string[i], "/")[len(strings.Split(files_string[i], "/")) - 1]
					split_j := strings.Split(files_string[j], "/")[len(strings.Split(files_string[j], "/")) - 1]
					i,err = strconv.Atoi(strings.Split(split_i, ".")[0])
					j,err =  strconv.Atoi(strings.Split(split_j, ".")[0])
 
					return i < j
				})
				break
			}

			var found_number bool
			for i:=0;i<len(image_files);i++{
				if strings.Split(file.Name(), ".")[0] == strconv.Itoa(i) {
					found_number = true
					break
				} 
			}
			if !found_number {
				Not_Numeric_Format_Checker += 1
			} else if found_number {
				Is_Numeric_Format_Checker += 1
			}
		}
		
		if len(image_files) != len(files_string) {
			fmt.Println("WARNING: Some images are missing")
			mp := make(map[string]struct{})
			for _,file := range files_string {
				mp[file] = struct{}{}
			}

			for _,file := range image_files {
				_,exists := mp[chapter_path + "/" + file.Name()]
				if exists == false {
					files_string = append(files_string, chapter_path + "/" + file.Name())
					}
				}
			mp = nil
			}
		image_files = nil
	} else {
		for _,file := range image_files {
			files_string = append(files_string, chapter_path + "/" + file.Name())
		}
	}

	var chapter_index int
	for i,j := range chapters {
		if j == chapter {
			chapter_index = i
			break
		}
	} 
	
	var next string
	var previous string
	if chapter_index != 0 {
		previous = chapters[chapter_index - 1]
	} else {
		previous = chapters[0]
	}

	if chapter_index != len(chapters) - 1 {
		next = chapters[chapter_index + 1]
	} else {
		next = chapters[chapter_index]
	}

	var chapters_links []string
	for _,i := range chapters {
		chapters_links = append(chapters_links, "/" + title + "/" + i)
	}

	d := manga_page_data{
		Images: files_string,
		Name: chapter,
		Next: "/" + title + "/" + next,
		Previous: "/" + title + "/" + previous,
		Chapters_names: chapters,
		Chapters_links: chapters_links,
	}


	tmpl,err := template.ParseFiles(filepath.Join("template", "manga_page.html"))
	if err != nil {
		fmt.Println("Something went wrong while parsing manga_page.html")
		fmt.Println(err)
		return
	}

	tmpl.Execute(w, d)


	file,err := os.ReadFile("data.json")
	if err != nil {
		os.Create("data.json")
	}
	
	var json_data []manga_information
	json.Unmarshal(file, &json_data)

	manga_index_in_json := findMangaIndexFromJson(json_data, title)
	
	if json_data[manga_index_in_json].Data.Last_read_chapter != chapter  {
		json_data[manga_index_in_json].Data.Last_read_chapter = chapter
		new_json_data,_ := json.MarshalIndent(json_data, "", "\t")
		os.WriteFile("data.json", new_json_data, 0600)
	}
}

func safeHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	var handler_already_made bool
	for _,handle := range handlers_made {
		if handle == pattern {
			handler_already_made = true
			break
		}
	}
	if !handler_already_made {
		http.HandleFunc(pattern, handler)
		handlers_made = append(handlers_made, pattern)
	}
}

func findMangaIndexFromJson(json_data []manga_information, title string) int {
	for i, j := range json_data {
		if j.Title == title {
			return i
		}
	}
	return -1
}
