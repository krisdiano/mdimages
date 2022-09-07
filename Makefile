MDPATH:=""

.PHONY: pull
pull:
	git pull

.PHONY: build
build:
	go build -o mdimage -gcflags='all=-N -l' main.go

.PHONY: dof
dof: pull _do

.PHONY: dos
dos: pull build _do


.PHONY: _do
_do:
	./mdimage  extract --path ${MDPATH} | awk '{printf("--paths %s\n", $$0)}' | tr '\n' ' ' | xargs  ./mdimage  upload | awk '{printf("--paths %s\n", $$0)}' | tr '\n' ' ' | xargs ./mdimage rewrite --path ${MDPATH} -i
