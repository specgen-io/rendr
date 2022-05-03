package render

import (
	"gopkg.in/specgen-io/yaml.v3"
	"io/ioutil"
	"path"
)

type Artifact struct {
	BuildCommand string `yaml:"build"`
}

const ArtifactFilename = ".rendr.yaml"

func GetArtifact(outPath string) (*Artifact, error) {
	artifactFullpath := path.Join(outPath, ArtifactFilename)
	data, err := ioutil.ReadFile(artifactFullpath)
	if err != nil {
		return nil, err
	}
	artifact := Artifact{}
	err = yaml.Unmarshal(data, &artifact)
	if err != nil {
		return nil, err
	}
	return &artifact, nil
}
