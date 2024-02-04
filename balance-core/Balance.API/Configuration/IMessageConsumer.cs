namespace Balance.API.Configuration;

public interface IMessageConsumer<out TValue> : IDisposable
{
    TValue? Consume(CancellationToken cancellationToken = default);
}