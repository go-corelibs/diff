#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

VERSION_TAGS += DIFF
DIFF_MK_SUMMARY := go-corelibs/diff
DIFF_MK_VERSION := v1.0.2

include CoreLibs.mk
