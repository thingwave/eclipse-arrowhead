SUBDIRS = ServiceRegistry Orchestrator Authorization DataManager

all:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir all; \
	done

all-arm64:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir all-arm64; \
	done

clean:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir clean; \
	done
