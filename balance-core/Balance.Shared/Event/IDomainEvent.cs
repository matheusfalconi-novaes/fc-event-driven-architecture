namespace Balance.Shared.Event;

public interface IDomainEvent
{
    public string Name { get; }
    public object Payload { get; }
    public DateTime DateTime { get; }
}


public abstract record DomainEvent(object Payload, DateTime DateTime) : IDomainEvent
{
    public abstract string Name { get; }
    
    public virtual object Payload { get; } = Payload;
    
    public DateTime DateTime { get; } = DateTime;
}