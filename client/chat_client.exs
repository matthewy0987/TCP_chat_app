defmodule ChatClient do
  def connect(host \\ "localhost", port \\ 6666) do
    {:ok, socket} = :gen_tcp.connect(host, port,
    [:binary, active: true, packet: :line])
    IO.puts("connected to #{host}:#{port}")
    start_input_loop(socket)
    message_loop(socket)
  end

  defp start_input_loop(socket) do
    spawn(fn -> input_loop(socket) end)
  end

  defp input_loop(socket) do
    case IO.gets("> ") do
      :eof ->
        IO.puts("DC")
        :gen_tcp.close(socket)

      data ->
        :gen_tcp.send(socket, data)
        input_loop(socket)
    end
  end

  defp message_loop(socket) do
    receive do
      {:tcp, ^socket, data} ->
        clear_line()
        IO.write(data)

        IO.write("> ")
        message_loop(socket)

      {:tcp_closed, ^socket} ->
        IO.puts("Server closed the connection.")
        :gen_tcp.close(socket)

      {:tcp_error, ^socket, reason} ->
        IO.puts("Error with socket: #{inspect(reason)}")
        :gen_tcp.close(socket)

      other ->
        IO.puts("Unexpected message: #{inspect(other)}")
        message_loop(socket)
    end
  end

  defp clear_line() do
    IO.write("\r\e[K")
  end

end

with [host, port | _] <- System.argv(),
  {port, ""} <- Integer.parse(port)
do
  ChatClient.connect(String.to_atom(host), port)
else
  _ -> ChatClient.connect(:localhost, 6666)
end
