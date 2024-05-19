package gitx

import (
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func GitClone(gitURL, path string) error {
	_, err := exec.LookPath("git")
	if err != nil {
		return err
	}
	c := exec.Command("git", "clone", gitURL)
	c.Dir = path
	c.Stderr = os.Stderr
	return c.Run()
}

func GitCheckout(branch, path string) error {
	if branch == "" {
		return nil
	}

	c := exec.Command("git", "checkout", branch)
	c.Dir = path
	c.Stderr = os.Stderr
	return c.Run()
}

func GitPath(gitURL string) (string, error) {
	if len(gitURL) > 3 && gitURL[0:3] == "git" {
		p := strings.Split(gitURL, "/")
		path := p[len(p)-1]
		return path[:len(path)-4], nil
	}
	u, err := url.Parse(gitURL)
	if err != nil {
		return "", err
	}
	p := strings.Split(strings.Trim(u.Path, ""), "/")
	path := p[len(p)-1]
	return path[:len(path)-4], nil
}
