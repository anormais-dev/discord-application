import { Monkey } from "../components/Monkey";

// Quem abre o link fora do Discord: a roleta é privada e só abre pela Activity.
export function OnlyDiscordPage() {
  return (
    <main className="private">
      <section className="panel private-card">
        <Monkey mood="dizzy" size={190} />
        <span className="tag tag-red private-tag">Acesso privado</span>
        <h1 className="px private-t">Só dentro do Discord</h1>
        <p className="private-d">
          A roleta é exclusiva dos Anormais e só abre pela <b>atividade</b>, dentro de um canal de voz.
        </p>
        <ol className="private-steps">
          <li>
            <span className="private-n">1</span>Entre num canal de voz do servidor
          </li>
          <li>
            <span className="private-n">2</span>Clique no foguete de atividades
          </li>
          <li>
            <span className="private-n">3</span>Escolha a Roleta dos Anormais
          </li>
        </ol>
      </section>
    </main>
  );
}
