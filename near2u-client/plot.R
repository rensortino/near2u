args = commandArgs(trailingOnly = TRUE)
my_data <- read.csv("./history.csv")
df <- data.frame(my_data)
print(args[1])
png("grafico.png")
df_sensor <- df[df$Code == args[1],]
plot(df$Value)
dev.off()