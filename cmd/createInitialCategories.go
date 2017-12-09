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

// createInitialCategoriesCmd represents the createInitialCategories command
var createInitialCategoriesCmd = &cobra.Command{
	Use:   "createInitialCategories",
	Short: "Creates initial categories",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createInitialCategories called")
		db, err := gorm.Open("postgres", "host=localhost user=perscorecal dbname=per_score_cal sslmode=disable password=perscorecal-dm")
		defer db.Close()
		if err != nil {
			log.Errorf("Error in createInitialCategories: %+v", err)
		}
		models.CreateInitialCategories(db)
	},
}

func init() {
	RootCmd.AddCommand(createInitialCategoriesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createInitialCategoriesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createInitialCategoriesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
