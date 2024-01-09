// local build srch file

package main

import (
	"fmt"
  "os"
  "errors"

	"srch"
  "github.com/urfave/cli/v2"
)

func main() {
  var recursive bool

  app := &cli.App{
    Name:  "srch",
    Usage: "Performant recursive file content search tool",
    Flags: []cli.Flag {
      &cli.BoolFlag {
        Name:         "recursive",
        Aliases:      []string{"r"},
        Value:        false,
        Usage:        "Flag for recursive directory searching",
        Destination:  &recursive,
      },
    },
    Action: func(cCtx *cli.Context) error {
      // Check number of args passed in before doing anything
      // 0      -> provide breaf usage explanation
      // not 2  -> return an error
      switch cCtx.Args().Len() {
      case 0:
        fmt.Println(`USAGE: srch [OPTION] PATTERN [DIRECTORY/FILE]
Try 'srch --help' for more information`)
        return nil
      case 2:
        break
      default:
        return errors.New(fmt.Sprintf("%s%d",
          "Expected 2 arguments... Got ", cCtx.Args().Len()))
      }

      // Get complete path to search
      path, err := srch.GetCompletePath(cCtx.Args().Get(1))
      if err != nil {
        return err
      }
      is_directory, err := srch.IsDirectory(path)
      if err != nil {
        return err
      }

      engine := srch.MakeStringFinder(cCtx.Args().Get(0))

      // If it is a file, simply call srch with the path
      if !is_directory {
        err := srch.Srch(engine, path)
        if err != nil {
          return err
        }
      } else {
        files, err := srch.GetFilesRecursively(path)
        if err != nil {
          return err
        }
        for _, e := range files {
          err := srch.Srch(engine, e)
          if err != nil {
            return err
          }
        }
      }

      return nil
    },
  }

  // Print any erorrs thrown
  if err := app.Run(os.Args); err != nil {
    fmt.Println(err.Error())
  }
}
