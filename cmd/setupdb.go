// Copyright © 2017 Vikram Anand <vikram.anand@renovite.com>
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

	"perScoreCal/models"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// setupdbCmd represents the setupdb command
var setupdbCmd = &cobra.Command{
	Use:   "setupdb",
	Short: "Setup DB",
	Long:  `Create and migrate schema`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setupdb called")
		db, err := gorm.Open("postgres", "host=localhost user=perscorecal dbname=per_score_cal sslmode=disable password=perscorecal-dm")
		defer db.Close()
		if err != nil {
			log.Errorf("Error in setupdb: %+v", err)
		}
		models.SetupDatabase(db)
	},
}

func init() {
	RootCmd.AddCommand(setupdbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupdbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupdbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
