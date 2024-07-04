// MIT License
//
// Copyright (c) 2024 chaunsin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package ncmctl

import (
	"context"
	"fmt"
	"os"

	"github.com/chaunsin/netease-cloud-music/api"
	"github.com/chaunsin/netease-cloud-music/api/weapi"
	"github.com/chaunsin/netease-cloud-music/pkg/log"

	"github.com/spf13/cobra"
)

type ExitOpts struct{}

type Exit struct {
	root *Root
	cmd  *cobra.Command
	opts ExitOpts
	l    *log.Logger
}

func NewExit(root *Root, l *log.Logger) *Exit {
	c := &Exit{
		root: root,
		l:    l,
		cmd: &cobra.Command{
			Use:     "exit",
			Short:   "Exit netease cloud music",
			Example: "  ncmctl exit",
		},
	}
	c.addFlags()
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.execute(cmd.Context(), args)
	}

	return c
}

func (c *Exit) addFlags() {}

func (c *Exit) Add(command ...*cobra.Command) {
	c.cmd.AddCommand(command...)
}

func (c *Exit) Command() *cobra.Command {
	return c.cmd
}

func (c *Exit) execute(ctx context.Context, args []string) error {
	cli, err := api.NewClient(c.root.Cfg.Network, c.l)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer cli.Close(ctx)
	request := weapi.New(cli)

	resp, err := request.Layout(ctx, &weapi.LayoutReq{})
	if err != nil {
		return fmt.Errorf("Layout: %w", err)
	}
	if resp.Code != 200 {
		return fmt.Errorf("Layout: %+v", resp)
	}

	// 只清理默认目录下得文件
	if err := os.Remove(c.root.Opts.Home + "/.ncmctl/cookie.json"); err != nil {
		log.Warn("remove cookie.json: %w", err)
	}
	c.cmd.Println("Logout success")
	return nil
}
