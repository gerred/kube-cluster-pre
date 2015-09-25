package virtualbox

type Option func(*Virtualbox)

func VagrantBox(url string) Option {
	return func(c *Virtualbox) {
		c.vagrantBox = url
	}
}

func VagrantBoxSHA1(sha1 string) Option {
	return func(c *Virtualbox) {
		c.vagrantBoxSHA1 = sha1
	}
}
