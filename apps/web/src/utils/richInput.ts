import type { JSONContent } from '@tiptap/core';

export function trimEmptyNodes(content: JSONContent) {
	if (!content.content || !Array.isArray(content.content)) {
		return content;
	}

	const cleaned = [...content.content];

	while (cleaned.length > 0 && isEmptyParagraph(cleaned[0])) {
		cleaned.shift();
	}

	while (cleaned.length > 0 && isEmptyParagraph(cleaned[cleaned.length - 1])) {
		cleaned.pop();
	}

	return { ...content, content: cleaned };
}

function isEmptyParagraph(node: any) {
	return (
		node.type === 'paragraph' &&
		(!node.content ||
			node.content.length === 0 ||
			(node.content.length === 1 && !node.content[0].text?.trim()))
	);
}

export function extractFirstNParagraphs(htmlString: string, n = 3) {
	const pattern = /<p>.*?<\/p>/gs;
	const matches = htmlString.match(pattern) || [];

	return { paragraphs: matches.slice(0, n).join(''), enoughMatches: matches.length >= n };
}
