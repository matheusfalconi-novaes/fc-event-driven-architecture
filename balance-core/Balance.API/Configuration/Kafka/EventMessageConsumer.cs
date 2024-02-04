using Balance.Domain.Event;
using Balance.Shared.Event;
using Confluent.Kafka;
using Confluent.Kafka.SyncOverAsync;
using Microsoft.Extensions.Options;

namespace Balance.API.Configuration.Kafka;

public record KafkaConfig
{
    public string BootStrapServers { get; set; } = null!;
    public string GroupId { get; set; } = null!;
    public string Topic { get; set; } = null!;
}

public abstract class EventMessageConsumer<TDomainEvent> : IMessageConsumer<TDomainEvent>
    where TDomainEvent : DomainEvent
{
    private readonly string _topic;
    private readonly IConsumer<string, TDomainEvent> _consumer;
    private readonly ILogger<EventMessageConsumer<TDomainEvent>> _logger;

    protected EventMessageConsumer(IOptions<KafkaConfig> configuration,
        ILogger<EventMessageConsumer<TDomainEvent>> logger)
    {
        _logger = logger;
        var consumerConfig = new ConsumerConfig
        {
            BootstrapServers = configuration.Value.BootStrapServers,
            GroupId = configuration.Value.GroupId,
            AutoOffsetReset = AutoOffsetReset.Earliest,
            EnableAutoCommit = true
        };
        _topic = configuration.Value.Topic;
        _consumer = new ConsumerBuilder<string, TDomainEvent>(consumerConfig)
            .SetKeyDeserializer(Deserializers.Utf8)
            .SetValueDeserializer(new JsonEventDeserializer<TDomainEvent>().AsSyncOverAsync())
            .Build();
        _consumer.Subscribe(_topic);
    }

    public TDomainEvent? Consume(CancellationToken cancellationToken = default)
    {
        if (_logger.IsEnabled(LogLevel.Information))
            _logger.LogInformation("Now consuming message for topic '{topic}' at: {time}.", _topic,
                DateTimeOffset.UtcNow);

        return _consumer.Consume(cancellationToken)?.Message?.Value;
    }

    public void Dispose()
    {
        _consumer.Close(); // Commit offsets and leave the group cleanly.
        _consumer.Dispose();
        GC.SuppressFinalize(this);
    }
}