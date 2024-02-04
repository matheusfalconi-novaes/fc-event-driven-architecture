using System.Text;
using Balance.Domain.Event;
using Balance.Shared.Event;
using Confluent.Kafka;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

namespace Balance.API.Configuration.Kafka;

internal class EventJsonConverter : JsonConverter
{
    public override bool CanConvert(Type objectType)
    {
        return objectType == typeof(DomainEvent);
    }

    public override void WriteJson(JsonWriter writer, object? value, JsonSerializer serializer)
    {
        throw new NotImplementedException();
    }

    public override object? ReadJson(JsonReader reader, Type objectType, object? existingValue,
        JsonSerializer serializer)
    {
        var jsonObject = JObject.Load(reader);
        return jsonObject["Name"]!.Value<string>() switch
        {
            BalanceUpdatedEvent.EventName => jsonObject.ToObject<BalanceUpdatedEvent>(serializer),
            _ => throw new NotImplementedException()
        };
    }

    public override bool CanWrite => false;
}

internal class JsonEventDeserializer<TDomainEvent> : IAsyncDeserializer<TDomainEvent> where TDomainEvent : DomainEvent
{
    private readonly JsonSerializerSettings _jsonSerializerSettings = new()
    {
        Converters = new List<JsonConverter> { new EventJsonConverter() }
    };

    public Task<TDomainEvent> DeserializeAsync(ReadOnlyMemory<byte> data, bool isNull, SerializationContext context)
    {
        var json = Encoding.UTF8.GetString(data.Span);
        var deserialized = JsonConvert.DeserializeObject<TDomainEvent>(json, _jsonSerializerSettings);
        if (deserialized is null)
            throw new InvalidDataException("Could not deserialize object");
        return Task.FromResult(deserialized);
    }
}