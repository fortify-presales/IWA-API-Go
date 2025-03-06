@echo on

sourceanalyzer "-Dcom.fortify.sca.ProjectRoot=.fortify" -b "iwa-api" -clean
sourceanalyzer "-Dcom.fortify.sca.ProjectRoot=.fortify" -b "iwa-api" "**/*.go" -verbose -debug
sourceanalyzer "-Dcom.fortify.sca.ProjectRoot=.fortify" -b "iwa-api" -verbose -debug -scan