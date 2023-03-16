# generate

To re-generate `xsd.go`, run:

    goxmlstruct -named-types -use-pointers-for-optional-fields=false xsd/*.xsd > internal/generate/xsd.go