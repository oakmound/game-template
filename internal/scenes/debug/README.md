# Debug
Here we expose a debug console that can be run by having the debug flag set.
For a ease of use just run mage dr to get a debug run in process.

Note that the current plumbing here is rather lossy as it takes a bit to realize that the main scene has changed contexts.
Given this expect this implementation to be deprecated in the near future (hopefully).