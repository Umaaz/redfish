// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package jobconfig

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("pickle.config.job#JobConfig", JobConfig{})
	pkl.RegisterMapping("pickle.config.job#Test", Test{})
	pkl.RegisterMapping("pickle.config.job#DataSource", DataSource{})
	pkl.RegisterMapping("pickle.config.job#RuleConfig", RuleConfig{})
	pkl.RegisterMapping("pickle.config.job#Expectation", Expectation{})
	pkl.RegisterMapping("pickle.config.job#JsonMatcher", JsonMatcherImpl{})
	pkl.RegisterMapping("pickle.config.job#StringMatcher", StringMatcherImpl{})
	pkl.RegisterMapping("pickle.config.job#FormBody", FormBodyImpl{})
	pkl.RegisterMapping("pickle.config.job#JsonBody", JsonBodyImpl{})
	pkl.RegisterMapping("pickle.config.job#JsonExtractor", JsonExtractorImpl{})
}
