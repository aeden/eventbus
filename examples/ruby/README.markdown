# Introduction

This is a Ruby event handler that can attach to the EventBus system and process events.

## Installation

If necessary, install bunder:

```
gem install bundler --no-ri --no-rdoc
```

Next, run bundler install:

```
bundle install
```

## Retry Logic

The current retry logic backs off connection retries for N * 1 second, where N is the retry count. It will stop retrying after a certain number of retries.
