package expansion

func Expand(src []string) (dst []string) {
	dst = src
	dst = expandShellParameters(dst)
	dst = expandFilenames(dst)
	return dst
}
