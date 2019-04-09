data <- read.csv(file=file.path("sensordata.csv"), header = TRUE, sep = ",", as.is = TRUE, strip.white = TRUE)

daniel <- data[data$Board=="30aea4505654",]
christian <- data[data$Board!="30aea4505654",]

daniel$Time <- as.POSIXct(daniel$Time / 10 ^ 9, origin = "1970-01-01")
christian$Time <- as.POSIXct(christian$Time / 10 ^ 9, origin = "1970-01-01")

plot(daniel$Time, daniel$Lightlevel, xlim = c(min(daniel$Time), max(daniel$Time)))
plot(christian$Time, christian$Lightlevel, xlim = c(min(christian$Time), max(christian$Time)))