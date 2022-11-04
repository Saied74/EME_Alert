The purpose of the application is to show which stations have been worked and
which ones have not during the ARRL EME contest.  The DVRA EME
station communications software is composed of LinRad for monitoring the band
at the full EME spectrum, MAP65 for decoding the traffic across the band and
WSJT-X for conducting Q95 QSOs.
The application builds a history of the QSOs from the WSJT-X log.  Immediately
after processing of the WSJT-X log file, it reads the MAP65 log file and compares
them to the history data.  It builds a structure of the logs not yet worked and
displays them on the screen.  The data displayed is:
+ Call sign
+ Time that Map65 logged the signal

The application is configured by the config.yaml file.  It contains four (or
  possibly five) lines.
+ The fully qualified path and the name of the WSJT-X log file.  
+ The fully qualified path and the filename of the MAP65 file.
+ The fully qualified path and filename of the test file
+ The start time (inclusive)
+ The end time (exclusive)

The purpose of the start and end times is to bracket the WSJT-X log data to
avoid QSOs that were conducted outside of the official contest time limits.

The configuration file must be located in $HOME/EME_Alert directory and must be
named config.yaml.

The test file mentioned above is used when the program is invoked with the -t
(test) flag.  In this case, a go routine is kicked off that reads the test file
and then every minute on the minute appends two line to the MAP65 log file and
one line to WSJT-X log file (the test starts with both files deleted or empty).
The algorithm is that in each minute represented by m, it adds lines 2m and 2m+1
to the Map65 log and line m to the WSJT-X log.  As a result, during each minute,
the length of display grows by one.
