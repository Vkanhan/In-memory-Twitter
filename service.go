package main

import (
	"errors"
	"sort"
	"sync"
	"time"
)

type TwitterService struct {
	users       map[int]*User
	tweets      map[int]*Tweet
	userTweets  map[int][]*Tweet
	lastUserID  int
	lastTweetID int
	mu          sync.RWMutex
}

func NewTwitterService() *TwitterService {
	return &TwitterService{
		users:      make(map[int]*User),
		tweets:     make(map[int]*Tweet),
		userTweets: make(map[int][]*Tweet),
	}
}

func (ts *TwitterService) CreateUser(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()

	for _, user := range ts.users {
		if user.Username == username {
			return nil, errors.New("username already exists")
		}
	}

	ts.lastUserID++
	newUser := &User{
		ID:        ts.lastUserID,
		Username:  username,
		Followers: make(map[int]bool),
		Following: make(map[int]bool),
	}

	ts.users[newUser.ID] = newUser
	ts.userTweets[newUser.ID] = []*Tweet{}
	return newUser, nil
}

func (ts *TwitterService) GetUser(userID int) (*User, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	user, exists := ts.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (ts *TwitterService) FollowUser(userID, targetUserID int) error {
	if userID == targetUserID {
		return errors.New("users cannot follow themselves")
	}

	user, err := ts.GetUser(userID)
	if err != nil {
		return err
	}

	targetUser, err := ts.GetUser(targetUserID)
	if err != nil {
		return err
	}

	user.mu.Lock()
	defer user.mu.Unlock()
	targetUser.mu.Lock()
	defer targetUser.mu.Unlock()

	if user.Following[targetUserID] {
		return errors.New("already following this user")
	}

	user.Following[targetUserID] = true
	targetUser.Followers[userID] = true
	return nil
}

func (ts *TwitterService) UnfollowUser(userID, targetUserID int) error {
	user, err := ts.GetUser(userID)
	if err != nil {
		return err
	}

	targetUser, err := ts.GetUser(targetUserID)
	if err != nil {
		return err
	}

	user.mu.Lock()
	defer user.mu.Unlock()
	targetUser.mu.Lock()
	defer targetUser.mu.Unlock()

	if !user.Following[targetUserID] {
		return errors.New("not following this user")
	}

	delete(user.Following, targetUserID)
	delete(targetUser.Followers, userID)
	return nil
}

func (ts *TwitterService) CreateTweet(userID int, content string) (*Tweet, error) {
	if content == "" {
		return nil, errors.New("tweet content cannot be empty")
	}

	_, err := ts.GetUser(userID)
	if err != nil {
		return nil, err
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.lastTweetID++
	newTweet := &Tweet{
		ID:        ts.lastTweetID,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now(),
		Likes:     0,
	}

	ts.tweets[newTweet.ID] = newTweet
	ts.userTweets[userID] = append(ts.userTweets[userID], newTweet)

	sort.Slice(ts.userTweets[userID], func(i, j int) bool {
		return ts.userTweets[userID][i].Timestamp.After(ts.userTweets[userID][j].Timestamp)
	})

	return newTweet, nil
}

func (ts *TwitterService) DeleteTweet(userID, tweetID int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	tweet, exists := ts.tweets[tweetID]
	if !exists {
		return errors.New("tweet not found")
	}

	if tweet.UserID != userID {
		return errors.New("unauthorized to delete this tweet")
	}

	delete(ts.tweets, tweetID)

	userTweets := ts.userTweets[userID]
	for i, t := range userTweets {
		if t.ID == tweetID {
			ts.userTweets[userID] = append(userTweets[:i], userTweets[i+1:]...)
			break
		}
	}

	return nil
}

func (ts *TwitterService) GetUserProfile(userID int) (map[string]any, error) {
	user, err := ts.GetUser(userID)
	if err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()
	user.mu.RLock()
	defer user.mu.RUnlock()

	followers := make([]int, 0, len(user.Followers))
	for followerID := range user.Followers {
		followers = append(followers, followerID)
	}

	following := make([]int, 0, len(user.Following))
	for followingID := range user.Following {
		following = append(following, followingID)
	}

	return map[string]any{
		"id":        user.ID,
		"username":  user.Username,
		"followers": followers,
		"following": following,
		"tweets":    ts.userTweets[userID],
	}, nil
}


func (ts *TwitterService) GetFeed(userID int, limit int) ([]*Tweet, error) {
	user, err := ts.GetUser(userID)
	if err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()
	user.mu.RLock()
	defer user.mu.RUnlock()

	if limit <= 0 {
		limit = 20
	}

	allTweets := make([]*Tweet, 0)
	for followedID := range user.Following {
		allTweets = append(allTweets, ts.userTweets[followedID]...)
	}
	allTweets = append(allTweets, ts.userTweets[userID]...)

	sort.Slice(allTweets, func(i, j int) bool {
		return allTweets[i].Timestamp.After(allTweets[j].Timestamp)
	})

	// Apply limit
	if limit > 0 && limit < len(allTweets) {
		allTweets = allTweets[:limit]
	}

	return allTweets, nil
}

func (ts *TwitterService) LikeTweet(tweetID int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	tweet, exists := ts.tweets[tweetID]
	if !exists {
		return errors.New("tweet not found")
	}

	tweet.Likes++
	return nil
}
