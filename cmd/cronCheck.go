// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

var k8s string

// cronCheckCmd represents the cronCheck command
var cronCheckCmd = &cobra.Command{
	Use:   "cronCheck",
	Short: "Check k8s update",
	Long:  `Use cronjob to check k8s new version and release note is update or not`,
	Run: func(cmd *cobra.Command, args []string) {
		cronjob()
	},
}

func init() {
	rootCmd.AddCommand(cronCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cronCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cronCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cronCheckCmd.Flags().StringVarP(&k8s, "k8s", "k", "v1.13", "k8s version will be updated")
}

func cronjob() {
	c := cron.New()
	c.AddFunc("0 * * * * *", func() { checkHomePage() })
	c.AddFunc("0 * * * * *", func() { checkReleaseNote() })
	c.Start()
	defer c.Stop()
	select {}
}

func checkHomePage() {
	fmt.Println("start to checkHomePage...")
	resp, err := http.Get("https://kubernetes.io")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if strings.Contains(string(body), k8s) {
		fmt.Println(Logo1)
		fmt.Printf("kubernetes %s HomePage Update!!!\n", k8s)
	}
	fmt.Println("home page is still not updated")
}

func checkReleaseNote() {
	fmt.Println("start to checkReleaseNote...")
	v := k8s[1:]
	resp, err := http.Get(fmt.Sprintf("https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-%s.md#major-themes", v))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if strings.Contains(string(body), "Major Themes") {
		fmt.Println(Logo2)
		fmt.Printf("kubernetes %s ReleaseNote Update!!!\n", k8s)
		return
	}
	fmt.Println("release note is still not updated")

}
