package main

import (
	"app/internal/config"
	"app/internal/db/postgres"
	"app/pkg/utils/helper"
	"log"
	"regexp"
	"strings"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

func main() {
	config.Parse("")
	postgres.Load()
	tables, err := postgres.Default().Migrator().GetTables()
	if err != nil {
		log.Fatal(err)
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:      "internal/db/postgres/query",
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		ModelPkgPath: "model",
		//FieldSignable: true,
	})
	g.UseDB(postgres.Default())
	re := regexp.MustCompile(`_\d+$`)
	tables = helper.ArrayFilter(tables, func(table string) bool {
		return !re.MatchString(table) && !strings.HasSuffix(table, "_bak")
	})
	for _, table := range tables {
		model := g.GenerateModel(table,
			gen.FieldTag("created_at", func(tag field.Tag) field.Tag {
				tag.Set("gorm", "column:created_at;not null;autoCreateTime")
				return tag
			}),
			gen.FieldTag("updated_at", func(tag field.Tag) field.Tag {
				tag.Set("gorm", "column:updated_at;not null;autoUpdateTime")
				return tag
			}),
		)
		g.ApplyBasic(model)
	}
	g.Execute()
}
