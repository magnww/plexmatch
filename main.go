package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"strconv"
	"strings"
)

var (
	cpuprofile = flag.String("cpu-profile", "", "write cpu profile to file")
	dir        = flag.String("dir", "", "")

	seRegexList     []*regexp.Regexp
	epRegexList     []*regexp.Regexp
	ignoreRegexList []*regexp.Regexp

	chineseNumbers = map[string]int64{
		"一": 1,
		"二": 2,
		"三": 3,
		"四": 4,
		"五": 5,
		"六": 6,
		"七": 7,
		"八": 8,
		"九": 9,
		"十": 10,
	}
)

func main() {
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ParseRegex()
	if *dir != "" {
		Detect(*dir)
	} else if os.Args[1] != "" {
		Detect(os.Args[1])
	}
}

func ParseRegex() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	workPath := filepath.Dir(ex)

	seRegexList = scanRegexpFile(path.Join(workPath, "session.txt"))
	epRegexList = scanRegexpFile(path.Join(workPath, "episode.txt"))
	ignoreRegexList = scanRegexpFile(path.Join(workPath, "ignore.txt"))
}

func scanRegexpFile(filepath string) []*regexp.Regexp {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	regexpList := make([]*regexp.Regexp, 0)
	for scanner.Scan() {
		regexpList = append(regexpList, regexp.MustCompile(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return regexpList
}

func Detect(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	builder := &strings.Builder{}

	title := removeNumber(path.Base(dir), epRegexList)
	title = removeIgnore(title, ignoreRegexList)
	log.Printf("title: %s", title)
	builder.WriteString(fmt.Sprintf("title: %s\n", title))

	se := match(dir, seRegexList)
	if se > 0 {
		log.Printf("se: %02d", se)
		builder.WriteString(fmt.Sprintf("season: %02d\n", se))
	}

	for _, file := range files {
		name := file.Name()
		ep := match(name, epRegexList)

		if ep > 0 {
			log.Printf("ep: %02d", ep)
			builder.WriteString(fmt.Sprintf("episode: %02d: %s\n", ep, name))
		}
	}

	f, err := os.Create(path.Join(dir, ".plexmatch"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(builder.String())
}

func match(name string, regexpList []*regexp.Regexp) int {
	var num int64
	var err error
	for _, regex := range regexpList {
		matches := regex.FindAllStringSubmatch(name, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				num, err = strconv.ParseInt(match[1], 10, 32)
				if err != nil {
					num = convertChineseNumber(match[1])
				}
				if num > 0 {
					break
				}
			}
		}
	}
	return int(num)
}

func removeNumber(name string, regexpList []*regexp.Regexp) string {
	for _, regex := range regexpList {
		matches := regex.FindAllStringSubmatch(name, -1)
		matchIndexes := regex.FindAllStringSubmatchIndex(name, -1)
		if len(matches) > 0 {
			for i, match := range matches {
				num, err := strconv.ParseInt(match[1], 10, 32)
				if err != nil {
					num = convertChineseNumber(match[1])
				}
				if num > 0 {
					return name[:matchIndexes[i][0]] + name[matchIndexes[i][1]:]
				}
			}
		}
	}
	return name
}

func removeIgnore(name string, regexpList []*regexp.Regexp) string {
	for _, regex := range regexpList {
		name = regex.ReplaceAllString(name, "")
	}
	return name
}

func convertChineseNumber(str string) int64 {
	return chineseNumbers[str]
}
