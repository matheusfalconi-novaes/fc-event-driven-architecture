using Balance.API.Configuration;
using Balance.API.Configuration.Kafka;
using Balance.API.Event.Handler;
using Balance.Domain.Event;
using Balance.Domain.Gateway;
using Balance.Domain.UseCase;
using Balance.Infrastructure.Gateway;
using Microsoft.AspNetCore.Mvc;

var builder = WebApplication.CreateBuilder(args);
builder.Configuration.AddEnvironmentVariables();

builder.Services.Configure<BalanceInfluxDbConfiguration>(builder.Configuration.GetSection("InfluxDB"));
builder.Services.Configure<KafkaConfig>(builder.Configuration.GetRequiredSection("KafkaConfig"));

builder.Services.AddSingleton<IBalanceGateway, BalanceInfluxDb>();
builder.Services.AddSingleton<FindAccountBalanceUseCase>();

builder.Services.AddSingleton<IMessageConsumer<BalanceUpdatedEvent>, BalanceUpdatedConsumer>();
builder.Services.AddHostedService<BalanceUpdatedService>();

var app = builder.Build();

app.MapGet("/balances/{accountId}",
    async (string accountId, [FromServices] FindAccountBalanceUseCase useCase) =>
    {
        var output = await useCase.Execute(new FindAccountBalanceInputDto(accountId));
        return output is not null ? Results.Ok(output) : Results.NotFound();
    });

Console.WriteLine("Server started");

app.Run();