package emb

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ListSkills returns the names of all skills in the embedded skills directory.
func ListSkills(embedFS fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(embedFS, "embed/skills")
	if err != nil {
		return nil, fmt.Errorf("read skills directory: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// CopySkills copies embedded skill files to the target project.
func CopySkills(embedFS fs.FS, targetDir string) error {
	skills, err := ListSkills(embedFS)
	if err != nil {
		return err
	}
	for _, name := range skills {
		src := filepath.Join("embed", "skills", name, "SKILL.md")
		dst := filepath.Join(targetDir, ".claude", "skills", name, "SKILL.md")
		if err := copyFile(embedFS, src, dst); err != nil {
			return fmt.Errorf("copy skill %s: %w", name, err)
		}
	}
	return nil
}

// CopySchema copies the myspec-driven schema to the target project.
func CopySchema(embedFS embed.FS, targetDir string) error {
	return copyDir(embedFS, "embed/schemas/myspec-driven", filepath.Join(targetDir, "openspec", "schemas", "myspec-driven"))
}

// RemoveSkills removes myspec skill files from the target project.
func RemoveSkills(embedFS fs.FS, targetDir string) error {
	skills, err := ListSkills(embedFS)
	if err != nil {
		return err
	}
	for _, name := range skills {
		os.RemoveAll(filepath.Join(targetDir, ".claude", "skills", name))
	}
	return nil
}

// RemoveSchema removes the myspec-driven schema from the target project.
func RemoveSchema(targetDir string) error {
	return os.RemoveAll(filepath.Join(targetDir, "openspec", "schemas", "myspec-driven"))
}

func copyFile(fsys fs.FS, src, dst string) error {
	in, err := fsys.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copyDir(fsys embed.FS, src, dst string) error {
	return fs.WalkDir(fsys, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(fsys, path, target)
	})
}
