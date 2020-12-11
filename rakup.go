package rakup

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rakyll/statik/fs"
	"github.com/sclevine/agouti"

	_ "github.com/dikmit/rakup/statik"
)

var (
	TOP_URL    = "https://websearch.rakuten.co.jp/"
	LOGIN_URL  = "https://grp03.id.rakuten.co.jp/rms/nid/login"
	SEARCH_URL = "https://websearch.rakuten.co.jp/Web?qt="
)

type Browser struct {
	driver *agouti.WebDriver
	page   *agouti.Page
}

func NewBrowser() (*Browser, error) {
	staticFS, err := fs.New()
	if err != nil {
		return nil, err
	}
	// crxBytes, err := ioutil.ReadFile("./rakuten.crx")
	r, err := staticFS.Open("/rakuten.crx")
	if err != nil {
		return nil, fmt.Errorf("static file read error: %w", err)
	}
	defer r.Close()
	crxBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			// "--user-data-dir=./profile",
			// "--headless",
			"--disable-gpu",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--window-size=1024,768",
			// "--disable-setuid-sandbox",
		}),
		agouti.ChromeOptions("prefs", map[string]string{
			"intl.accept_languages": "ja,en_US",
		}),
		agouti.ChromeOptions("extensions", [][]byte{crxBytes}),
		// agouti.Debug,
	)

	if err := driver.Start(); err != nil {
		return nil, err
	}
	page, err := driver.NewPage()
	if err != nil {
		if err := driver.Stop(); err != nil {
			return nil, fmt.Errorf("driver Stop error: %w", err)
		}
		return nil, err
	}
	return &Browser{driver: driver, page: page}, nil
}

func (b *Browser) Close() error {
	if err := b.driver.Stop(); err != nil {
		return fmt.Errorf("driver stop error: %w\n", err)
	}
	return nil
}

func (b *Browser) Login(user string, password string) error {
	// Login
	if err := b.page.Navigate(LOGIN_URL); err != nil {
		return err
	}
	if err := b.page.FindByID("loginInner_u").Fill(user); err != nil {
		return err
	}
	if err := b.page.FindByID("loginInner_p").Fill(password); err != nil {
		return err
	}
	if err := b.page.Find("input.loginButton").Click(); err != nil {
		return err
	}
	if err := b.page.Find("#contents > form > p.submit > input[type=submit]").Click(); err != nil {
		return fmt.Errorf("login error")
	}
	return nil
}

// GetTrend is トレンドワードを取得
func (b *Browser) GetTrend(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	errCh := make(chan error)
	ch := make(chan []string)

	go func() {
		// 検索ページTOP
		if err := b.page.Navigate(TOP_URL); err != nil {
			errCh <- err
			return
		}
		html, err := b.page.HTML()
		if err != nil {
			errCh <- err
			return
		}
		r := strings.NewReader(html)
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			errCh <- err
			return
		}
		todayWords := []string{}
		doc.Find("#simple-top-search-form>div.sc-AxhCb>dl>dd>div").Each(func(i int, s *goquery.Selection) {
			word, ok := s.Find("input").Attr("name")
			if ok {
				todayWords = append(todayWords, word)
			}
		})
		ch <- todayWords
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-errCh:
		return nil, e
	case words := <-ch:
		return words, nil
	}
}

// Search is do search with word. return kuchi int, error
func (b *Browser) Search(word string) (int, error) {
	if err := b.page.Navigate(SEARCH_URL + word); err != nil {
		return 0, err
	}
	// 口数がページ上でレンダリングされるのを待つ
	time.Sleep(1 * time.Second)
	s, err := b.page.First("em").Text()
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("kuchisu Atoi error: %w", err)
	}
	return i, nil
}

func ReadWords(file string) (words []string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		words = append(words, s.Text())
	}
	if err := s.Err(); err != nil {
		return words, err
	}
	return words, nil
}

func SaveWord(s string, file string) (err error) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	if _, err := f.WriteString(s + "\n"); err != nil {
		return err
	}
	return nil
}
