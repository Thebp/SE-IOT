data = read.csv('~/Documents/Github/SE-IOT/Assignments/Assignment1-full-circle-latency/Devices/PC/data.txt', header = FALSE, sep = ",")

# Add a Normal Curve (Thanks to Peter Dalgaard)
x <- data$V3

y <- subset(x, data$V3>0.35)

h<-hist(x, breaks=200, col="gray", xlab="Latency in seconds",
        main="Full circle latency - 1000 samples")
xfit<-seq(min(x),max(x),length=100)
yfit<-dnorm(xfit,mean=mean(x),sd=sd(x))
yfit <- yfit*diff(h$mids[1:2])*length(x)
lines(xfit, yfit, col="blue", lwd=2)


# Kernel Density Plot
d <- density(x) # returns the density data
plot(d) # plots the results 

# Filled Density Plot
d <- density(x)
plot(d, main="test")
polygon(d, col="red", border="blue")

# Rainbow Barplot
r <- barplot(x, col=rainbow(100))

library(lattice)
densityplot(~ V3, data = data)


data100 = read.csv('~/Documents/Github/SE-IOT/Assignments/Assignment1-full-circle-latency/Devices/PC/data100.txt', header = FALSE, sep = ",")

# Add a Normal Curve (Thanks to Peter Dalgaard)
x100 <- data100$V3

