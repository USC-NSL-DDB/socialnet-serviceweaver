package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"SocialNetwork/shared/api"
	"SocialNetwork/shared/common"
)

// var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var r = rand.New(rand.NewSource(42))

func randomString(letters string, length int) string {
	var result strings.Builder
	result.Grow(length) // Pre-allocate memory for efficiency
	for i := 0; i < length; i++ {
		result.WriteByte(letters[r.Intn(len(letters))])
	}
	return result.String()
}

// randomDigits generates a random string of digits with a specified length
func randomDigits(length int) int64 {
	const digits = "0123456789"
	var sb strings.Builder
	sb.Grow(length) // Pre-allocate memory to improve performance

	for i := 0; i < length; i++ {
		sb.WriteByte(digits[r.Intn(len(digits))])
	}

	result, err := strconv.ParseInt(sb.String(), 10, 64)
	if err != nil {
		fmt.Println("Error parsing result:", err)
		return 0
	}
	return result
}

// User registration
func registerUser(addr string, user string, wg *sync.WaitGroup) {
	defer wg.Done()

	user_id, err := strconv.ParseInt(user, 10, 64)
	if err != nil {
		fmt.Println("Error parsing user ID:", err)
		return
	}

	req := &api.RegisterUserWithIdRequest{
		FirstName: "first_name_" + user,
		LastName:  "last_name_" + user,
		Username:  "username_" + user,
		Password:  "password_" + user,
		UserId:    user_id,
	}

	api.RegisterUserWithId(addr, req)
}

// User following
func followUser(addr string, followerID, followeeID string, wg *sync.WaitGroup) {
	defer wg.Done()

	req := &api.FollowWithUsernameRequest{
		Username:         "username_" + followerID,
		FolloweeUsername: "username_" + followeeID,
	}
	api.FollowWithUsername(addr, req)
}

// Compose post
func composePost(addr string, user_id int, num_users int, wg *sync.WaitGroup) {
	defer wg.Done()

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	text := randomString(letters, 256)

	// User mentions
	numMentions := r.Intn(6)
	for i := 0; i < numMentions; i++ {
		text += " @username_" + strconv.Itoa(r.Intn(num_users)+1)
	}

	// URLs
	numURLs := r.Intn(6)
	for i := 0; i < numURLs; i++ {
		text += " http://" + randomString("abcdefghijklmnopqrstuvwxyz0123456789", 64)
	}

	// Media
	mediaIDs := make([]int64, 0)
	mediaTypes := make([]string, 0)
	numMedia := r.Intn(6)
	for i := 0; i < numMedia; i++ {
		mediaIDs = append(mediaIDs, randomDigits(18))
		mediaTypes = append(mediaTypes, "png")
	}

	req := &api.ComposePostRequest{
		UserId:     int64(user_id),
		Username:   "username_" + strconv.Itoa(user_id),
		Text:       text,
		MediaIds:   mediaIDs,
		MediaTypes: mediaTypes,
		PostType:   common.POST,
	}

	api.ComposePost(addr, req)
}

// getNodes reads the first line from the given file and converts the first word to an integer
func getNodes(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() { // Read the first line
		line := scanner.Text()
		words := strings.Fields(line)
		if len(words) > 0 {
			val, err := strconv.Atoi(words[0]) // Convert the first word to an integer
			if err != nil {
				return 0
			} else {
				return val
			}
		}
		return 0
	}
	if err := scanner.Err(); err != nil {
		return 0
	}

	return 0
}

func getEdges(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	strconv.Atoi(scanner.Text())
	var edges [][]string

	for scanner.Scan() {
		edges = append(edges, strings.Fields(scanner.Text()))
	}

	return edges
}

func register(addr string, nodes int) {
	idx := 0
	fmt.Println("Registering Users...")
	var wg sync.WaitGroup
	for i := 1; i <= nodes+1; i++ {
		idx += 1
		wg.Add(1)
		go registerUser(addr, strconv.Itoa(i), &wg)
		if idx%100 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func follow(addr string, edges [][]string) {
	idx := 0
	fmt.Println("Adding follows...")
	var wg sync.WaitGroup
	for _, edge := range edges {
		idx += 1
		wg.Add(2)
		go followUser(addr, edge[0], edge[1], &wg)
		go followUser(addr, edge[1], edge[0], &wg)
		// wg.Wait()
		// return
		if idx%50 == 0 {
			wg.Wait()
			fmt.Println("Added", idx*2)
			// return
			time.Sleep(1 * time.Second)
		}
	}
	wg.Wait()
}

func compose(addr string, nodes int) {
	idx := 0
	fmt.Println("Composing posts...")
	var wg sync.WaitGroup
	for i := 0; i <= nodes; i++ {
		upper := r.Intn(20)
		for j := 0; j < upper; j++ {
			idx += 1
			wg.Add(1)
			go composePost(addr, i, nodes, &wg)
			if idx%100 == 0 {
				wg.Wait()
			}
		}
	}
	wg.Wait()
}

func main() {
	// Command line arguments and initialization logic here similar to Python code
	addr := api.BASE_URL
	filepath := "./social-graph/socfb-Reed98/socfb-Reed98.mtx"
	nodes := getNodes(filepath)
	edges := getEdges(filepath)

	fmt.Println("addr:", addr)
	fmt.Println("Nodes:", nodes)
	fmt.Println("Edges:", len(edges))

	fmt.Println("First 10 edges:")
	for i := 0; i < 10 && i < len(edges); i++ {
		fmt.Println(edges[i])
	}

	register(addr, nodes)
	follow(addr, edges)
	// compose(addr, nodes)
}
