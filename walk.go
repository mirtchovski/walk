// Copyright 2010 Andrey Mirtchovski. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the Go distribution's
// LICENSE file.

// A non-recursive filesystem walker
package walk

import (
	"fmt"
	"os"
	"path/filepath"
)

type queue struct {
	name string
	seen bool
	cd   bool
	next *queue
}

var lstat = os.Lstat // for testing

// Walkiter iteratively descends through a directory storing subdirectories
// on s and calling walkFn for each file or directory it encounters
func walkiter(s *queue, walkFn filepath.WalkFunc) (haderror error) {
	for {
		if s == nil {
			return haderror
		}
		if s.seen {
			if s.cd {
				err := os.Chdir("..")
				if err != nil {
					cwd, err2 := os.Getwd()
					return fmt.Errorf("can't dotdot (..); : %s from %s: %v (getwd error: %v)", s.name, cwd, err, err2)
				}
			}
			s = s.next
			continue
		}
		s.seen = true

		ourname := s.name

		info, err := lstat(filepath.Base(ourname))
		if err != nil {
			haderror = walkFn(ourname, info, err)
			continue
		}

		err = walkFn(ourname, info, err)
		if err != nil {
			if info.IsDir() && err == filepath.SkipDir {
				continue
			}
			haderror = err
			continue
		}

		// if a directory, chdir and list it
		if info.IsDir() {
			err := os.Chdir(filepath.Base(ourname))
			if err != nil {
				haderror = walkFn(ourname, info, err)
				continue
			}
			s.cd = true
			file, err := os.Open(".")
			if err != nil {
				haderror = walkFn(ourname, info, err)
				continue
			}
			names, err := file.Readdirnames(0)
			if err != nil {
				haderror = walkFn(ourname, info, err)
			}
			file.Close()
			for i := len(names) - 1; i >= 0; i-- {
				ns := new(queue)
				ns.name = filepath.Join(ourname, names[i])
				ns.next = s
				s = ns
			}
		}
	}
	panic("unreachable")
}

// Walk does a non-recursive walk of the directory rooted at path, calling f for
// each file it encounters. The walk will descend using Chdir, so that deeply nested
// paths longer than PATH_MAX (1024 on OSX) would still be reachable.
func Walk(root string, walkFn filepath.WalkFunc) error {
	dir, err := os.Getwd()
	if err != nil {
		return err // wouldn't want to leave caller in an unknown dir
	}
	defer os.Chdir(dir)

	err = os.Chdir(filepath.Dir(root))
	if err != nil {
		return walkFn(root, nil, err)
	}
	s := new(queue)
	s.name = root
	return walkiter(s, walkFn)
}
