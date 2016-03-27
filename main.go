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
)

func rankFilesByOccurance(wordFrequencies map[string]int) PairList{
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key string
	Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }


func main (){
	var occurancesList = make(map[string]int)

	searchDir := "./"
	csprojRegex := regexp.MustCompile(".csproj")
	fileIncludeRegex := regexp.MustCompile("(<Content Include=\")(.*)\"")


	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		isProjFile := csprojRegex.FindString(path)
		if isProjFile != ""{
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {

	}



	for _, file := range fileList {
		fileReader, _ := os.Open(file)
		scanner := bufio.NewScanner(fileReader)
		for scanner.Scan() {
			lineOfFile := scanner.Text()
			jumpOverLine := fileIncludeRegex.FindAllStringSubmatch(lineOfFile, -2)
			//Process line as part of migration
			if len(jumpOverLine) != 0 {
				occurancesList[jumpOverLine[0][2]] += 1
			}
		}

	}

	pairList := rankFilesByOccurance(occurancesList)
	for _, occurance := range pairList {
		fmt.Println(occurance)

	}
}