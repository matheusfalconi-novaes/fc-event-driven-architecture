using Balance.Shared.Event;
using Newtonsoft.Json;

namespace Balance.Domain.Event;

public record BalanceUpdatedEventPayload
{
    [JsonProperty("account_id_from")] public string AccountIdFrom { get; set; } = null!;
    [JsonProperty("account_id_to")] public string AccountIdTo { get; set; } = null!;
    [JsonProperty("balance_account_from")] public float BalanceAccountFrom { get; set; }
    [JsonProperty("balance_account_to")] public float BalanceAccountTo { get; set; }
}

public record BalanceUpdatedEvent(BalanceUpdatedEventPayload Payload, DateTime DateTime) : DomainEvent(Payload, DateTime)
{
    public const string EventName = "BalanceUpdated";
    public override string Name => EventName;
    public new BalanceUpdatedEventPayload Payload { get; } = Payload;

    public BalanceUpdatedEvent() : this(new BalanceUpdatedEventPayload(), DateTime.UtcNow)
    {
    }
}