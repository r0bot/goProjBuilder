package main

import (
	//"io/ioutil"
	"fmt"
	//"log"
	"path/filepath"
	"os"
	"regexp"
	"bufio"
	"sort"
"strings"
)

//Sorting
func rankFilesByOccurance(wordFrequencies map[string][]string) PairList{
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, len(v), v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key string
	Value int
	Payload []string
}

type PairList []Pair
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

//Sorting End

func main (){
	//Init output file.
	structureFile, _ := os.Create("./structure.txt")
	defer structureFile.Close()

	//Init directory and regexs
	searchDir := "../../../backend-services/Server/APIServer"
	csprojRegex := regexp.MustCompile(".csproj")
	fileIncludeRegex := regexp.MustCompile("(<Content Include=\")(.*)\"")

	//Init the maps needed to hold the data
	occurancesList := make(map[string][]string)
	fileList := make(map[string]os.FileInfo)

	//Walk the root directory and look for .csproj files
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		isProjFile := csprojRegex.FindString(path)
		if isProjFile != ""{
			fileList[path] = f
		}
		return nil
	})

	if err != nil {
		//TODO handle the error
	}

	//Scan each csproj file and extract the includes
	for path, fileInfo := range fileList {
		fileReader, _ := os.Open(path)
		scanner := bufio.NewScanner(fileReader)
		for scanner.Scan() {
			lineOfFile := scanner.Text()
			jumpOverLine := fileIncludeRegex.FindAllStringSubmatch(lineOfFile, -2)
			//Process line as part of migration
			if len(jumpOverLine) != 0 {
				occurancesList[jumpOverLine[0][2]] = append(occurancesList[jumpOverLine[0][2]], fileInfo.Name())
			}
		}

	}

	//Order the extracted includes based on the number of times they appear in different csproj files.
	pairList := rankFilesByOccurance(occurancesList)

	//Output the results to a file
	for _, occurance := range pairList {
		fmt.Println(occurance.Payload)
		structureFile.WriteString(string(occurance.Key) + " : " + string(occurance.Value) + "\n")
		structureFile.WriteString(strings.Join(occurance.Payload, " ") + "\n")
		structureFile.WriteString("----------------------------------------------------------------\n\n\n")
	}
}