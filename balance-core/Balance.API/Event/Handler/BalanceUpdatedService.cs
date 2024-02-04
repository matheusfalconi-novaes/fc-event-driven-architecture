using Balance.API.Configuration;
using Balance.API.Configuration.Kafka;
using Balance.Domain.Event;
using Balance.Domain.Gateway;
using Microsoft.Extensions.Options;

namespace Balance.API.Event.Handler;

internal class BalanceUpdatedConsumer(IOptions<KafkaConfig> configuration, ILogger<BalanceUpdatedConsumer> logger)
    : EventMessageConsumer<BalanceUpdatedEvent>(configuration, logger)
{
    public new BalanceUpdatedEvent? Consume(CancellationToken cancellationToken = default)
        => base.Consume(cancellationToken);
}

public class BalanceUpdatedService(
    IMessageConsumer<BalanceUpdatedEvent> messageConsumer,
    IBalanceGateway balanceGateway,
    ILogger<BalanceUpdatedService> logger) : BackgroundService
{
    protected override Task ExecuteAsync(CancellationToken stoppingToken)
        => Task.Run(() => ProcessBalance(stoppingToken), stoppingToken);

    private async Task ProcessBalance(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            var @event = messageConsumer.Consume(stoppingToken);
            if (@event is null)
            {
                logger.LogWarning("Message contains null value.");
                continue;
            }

            var eventPayload = @event.Payload;
            var balances = new[]
            {
                new BalanceData(eventPayload.AccountIdFrom, eventPayload.BalanceAccountFrom, @event.DateTime),
                new BalanceData(eventPayload.AccountIdTo, eventPayload.BalanceAccountTo, @event.DateTime)
            };
            await balanceGateway.SaveBalances(balances, stoppingToken);
        }

        logger.LogDebug($"{nameof(BalanceUpdatedService)} background task is stopping.");
    }

    public override void Dispose()
    {
        messageConsumer.Dispose();
        base.Dispose();
        GC.SuppressFinalize(this);
    }
}